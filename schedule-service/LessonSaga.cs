using System.Net.Http.Json;
using Npgsql;
namespace University;

public class LessonSaga(LessonRepository repository, IHttpClientFactory factory,
    IConfiguration config, ILogger<LessonSaga> logger)
{
    public async Task<SagaResult> ExecuteAsync(LessonCommand command, CancellationToken ct)
    {
        if (command.TeacherId == Guid.Empty || command.SubjectId == Guid.Empty || command.RoomId == Guid.Empty
            || string.IsNullOrWhiteSpace(command.Group) || command.EndsAt <= command.StartsAt)
            throw new ArgumentException("Invalid lesson");
        var id = Guid.NewGuid();
        var http = factory.CreateClient();
        var roomUrl = config["RoomsUrl"]!;
        logger.LogInformation("SAGA {Id} STARTED", id);
        var bookingAttempted = false;
        try
        {
            using var teacherResponse = await http.GetAsync($"{config["TeachersUrl"]}/api/teachers/{command.TeacherId}", ct);
            teacherResponse.EnsureSuccessStatusCode();
            var students = await http.GetFromJsonAsync<List<Student>>($"{config["StudentsUrl"]}/api/students", ct);
            var count = students!.Count(s => s.Group == command.Group);
            if (count == 0) throw new InvalidOperationException("Student group not found");
            await using (var db = await repository.OpenAsync(ct))
            await using (var query = new NpgsqlCommand("SELECT EXISTS(SELECT 1 FROM subjects WHERE id=@id)", db))
            {
                query.Parameters.AddWithValue("id", command.SubjectId);
                if (!(bool)(await query.ExecuteScalarAsync(ct))!) throw new InvalidOperationException("Subject not found");
            }
            bookingAttempted = true;
            using var response = await http.PostAsJsonAsync($"{roomUrl}/api/bookings", new
            {
                id, roomId = command.RoomId, startsAt = command.StartsAt,
                endsAt = command.EndsAt, size = count
            }, ct);
            response.EnsureSuccessStatusCode();
            logger.LogInformation("SAGA {Id} Room reserved", id);
            if (command.FailAfterBooking) throw new InvalidOperationException("Test failure after booking");
            await repository.SaveAsync(new(id, command.TeacherId, command.SubjectId, command.RoomId,
                command.Group, command.StartsAt, command.EndsAt, "Confirmed"), ct);
            logger.LogInformation("SAGA {Id} COMPLETED", id);
            return new(true, id, "Lesson created", false);
        }
        catch (Exception e)
        {
            var compensated = true;
            logger.LogWarning("SAGA {Id} COMPENSATION STARTED: {Message}", id, e.Message);
            if (bookingAttempted)
            {
                try
                {
                    using var cancellation = new CancellationTokenSource(TimeSpan.FromSeconds(5));
                    using var response = await http.DeleteAsync($"{roomUrl}/api/bookings/{id}", cancellation.Token);
                    response.EnsureSuccessStatusCode();
                }
                catch (Exception error)
                {
                    compensated = false;
                    logger.LogError("SAGA {Id} Booking cancellation failed: {Message}", id, error.Message);
                }
            }
            try
            {
                await using var db = await repository.OpenAsync();
                await using var delete = new NpgsqlCommand("DELETE FROM lessons WHERE id=@id", db);
                delete.Parameters.AddWithValue("id", id);
                await delete.ExecuteNonQueryAsync();
            }
            catch (Exception error)
            {
                compensated = false;
                logger.LogError("SAGA {Id} Lesson deletion failed: {Message}", id, error.Message);
            }
            logger.LogWarning("SAGA {Id} COMPENSATION {Status}", id, compensated ? "COMPLETED" : "INCOMPLETE");
            return new(false, id, e.Message, compensated);
        }
    }
}
