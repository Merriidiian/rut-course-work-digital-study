using Npgsql;
namespace University;

public class LessonRepository(IConfiguration config)
{
    public async Task<NpgsqlConnection> OpenAsync(CancellationToken ct = default)
    {
        var db = new NpgsqlConnection(config.GetConnectionString("Postgres"));
        await db.OpenAsync(ct);
        return db;
    }

    public async Task SaveAsync(Lesson lesson, CancellationToken ct = default)
    {
        await using var db = await OpenAsync(ct);
        await using var command = new NpgsqlCommand("""
            INSERT INTO lessons VALUES (@id,@teacher,@subject,@room,@group,@start,@end,@status)
            """, db);
        command.Parameters.AddWithValue("id", lesson.Id);
        command.Parameters.AddWithValue("teacher", lesson.TeacherId);
        command.Parameters.AddWithValue("subject", lesson.SubjectId);
        command.Parameters.AddWithValue("room", lesson.RoomId);
        command.Parameters.AddWithValue("group", lesson.Group);
        command.Parameters.AddWithValue("start", lesson.StartsAt.ToUniversalTime());
        command.Parameters.AddWithValue("end", lesson.EndsAt.ToUniversalTime());
        command.Parameters.AddWithValue("status", lesson.Status);
        await command.ExecuteNonQueryAsync(ct);
    }

    public async Task<List<Lesson>> GetAsync(string? group, Guid? teacher, CancellationToken ct)
    {
        await using var db = await OpenAsync(ct);
        await using var command = new NpgsqlCommand("""
            SELECT * FROM lessons WHERE (@group='' OR group_name=@group)
              AND (@teacher='00000000-0000-0000-0000-000000000000'::uuid OR teacher_id=@teacher)
            ORDER BY starts_at
            """, db);
        command.Parameters.AddWithValue("group", group ?? "");
        command.Parameters.AddWithValue("teacher", teacher ?? Guid.Empty);
        await using var reader = await command.ExecuteReaderAsync(ct);
        var result = new List<Lesson>();
        while (await reader.ReadAsync(ct))
            result.Add(new(reader.GetGuid(0), reader.GetGuid(1), reader.GetGuid(2), reader.GetGuid(3),
                reader.GetString(4), reader.GetFieldValue<DateTimeOffset>(5),
                reader.GetFieldValue<DateTimeOffset>(6), reader.GetString(7)));
        return result;
    }
}
