using Npgsql;
using University.Core.Clients;
using University.Contracts;
using University.Core.Repositories;

namespace University.Core.Endpoints;

public static class SubjectEndpoints
{
    public static void MapSubjectEndpoints(this WebApplication app)
    {
        var routes = app.MapGroup("/api/subjects");
        routes.MapGet("", async (
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
            Results.Ok(await repository.GetSubjectsAsync(cancellationToken)));

        routes.MapGet("/{id:guid}", async (
            Guid id,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            var result = await repository.FindSubjectAsync(id, cancellationToken);
            return result is null ? Results.NotFound(new ApiError("Record not found")) : Results.Ok(result);
        });

        routes.MapPost("", async (
            SubjectRequest request,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            if (string.IsNullOrWhiteSpace(request.Name) || request.Hours <= 0)
            {
                return Results.BadRequest(new ApiError("Invalid subjects"));
            }

            var result = await repository.SaveSubjectAsync(null, request, cancellationToken);
            return Results.Created($"/api/subjects/{result!.Id}", result);
        });

        routes.MapPut("/{id:guid}", async (
            Guid id,
            SubjectRequest request,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            if (string.IsNullOrWhiteSpace(request.Name) || request.Hours <= 0)
            {
                return Results.BadRequest(new ApiError("Invalid subjects"));
            }

            var result = await repository.SaveSubjectAsync(id, request, cancellationToken);
            return result is null ? Results.NotFound(new ApiError("Record not found")) : Results.Ok(result);
        });

        routes.MapDelete("/{id:guid}", async (
            Guid id,
            CatalogRepository repository,
            ScheduleClient scheduleClient,
            CancellationToken cancellationToken) =>
        {
            try
            {
                if (await scheduleClient.HasSubjectLessonsAsync(id, cancellationToken))
                {
                    return Results.Conflict(new ApiError("Record is used by a lesson"));
                }

                var deleted = await repository.DeleteSubjectAsync(id, cancellationToken);
                return deleted ? Results.NoContent() : Results.NotFound(new ApiError("Record not found"));
            }
            catch (HttpRequestException)
            {
                return Results.Json(new ApiError("Schedule service is unavailable"), statusCode: 502);
            }
            catch (TaskCanceledException)
            {
                return Results.Json(new ApiError("Schedule service timed out"), statusCode: 502);
            }
            catch (PostgresException exception) when (exception.SqlState == "23503")
            {
                return Results.Conflict(new ApiError("Record is used by a lesson"));
            }
        });
    }
}
