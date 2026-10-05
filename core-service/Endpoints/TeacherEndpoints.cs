using Npgsql;
using University.Core.Models;
using University.Core.Repositories;

namespace University.Core.Endpoints;

public static class TeacherEndpoints
{
    public static void MapTeacherEndpoints(this WebApplication app)
    {
        var routes = app.MapGroup("/api/teachers");
        routes.MapGet("", async (
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
            Results.Ok(await repository.GetTeachersAsync(cancellationToken)));

        routes.MapGet("/{id:guid}", async (
            Guid id,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            var result = await repository.FindTeacherAsync(id, cancellationToken);
            return result is null ? Results.NotFound() : Results.Ok(result);
        });

        routes.MapPost("", async (
            TeacherRequest request,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            if (string.IsNullOrWhiteSpace(request.Name))
            {
                return Results.BadRequest("Invalid teachers");
            }

            var result = await repository.SaveTeacherAsync(null, request, cancellationToken);
            return Results.Created($"/api/teachers/{result!.Id}", result);
        });

        routes.MapPut("/{id:guid}", async (
            Guid id,
            TeacherRequest request,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            if (string.IsNullOrWhiteSpace(request.Name))
            {
                return Results.BadRequest("Invalid teachers");
            }

            var result = await repository.SaveTeacherAsync(id, request, cancellationToken);
            return result is null ? Results.NotFound() : Results.Ok(result);
        });

        routes.MapDelete("/{id:guid}", async (
            Guid id,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            try
            {
                var deleted = await repository.DeleteTeacherAsync(id, cancellationToken);
                return deleted ? Results.NoContent() : Results.NotFound();
            }
            catch (PostgresException exception) when (exception.SqlState == "23503")
            {
                return Results.Conflict("Record is used by a lesson");
            }
        });
    }
}
