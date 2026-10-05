using Npgsql;
using University.Core.Models;
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
            return result is null ? Results.NotFound() : Results.Ok(result);
        });

        routes.MapPost("", async (
            SubjectRequest request,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            if (string.IsNullOrWhiteSpace(request.Name) || request.Hours <= 0)
            {
                return Results.BadRequest("Invalid subjects");
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
                return Results.BadRequest("Invalid subjects");
            }

            var result = await repository.SaveSubjectAsync(id, request, cancellationToken);
            return result is null ? Results.NotFound() : Results.Ok(result);
        });

        routes.MapDelete("/{id:guid}", async (
            Guid id,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            try
            {
                var deleted = await repository.DeleteSubjectAsync(id, cancellationToken);
                return deleted ? Results.NoContent() : Results.NotFound();
            }
            catch (PostgresException exception) when (exception.SqlState == "23503")
            {
                return Results.Conflict("Record is used by a lesson");
            }
        });
    }
}
