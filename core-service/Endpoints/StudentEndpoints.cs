using Npgsql;
using University.Contracts;
using University.Core.Repositories;

namespace University.Core.Endpoints;

public static class StudentEndpoints
{
    public static void MapStudentEndpoints(this WebApplication app)
    {
        var routes = app.MapGroup("/api/students");
        routes.MapGet("", async (
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
            Results.Ok(await repository.GetStudentsAsync(cancellationToken)));

        routes.MapGet("/{id:guid}", async (
            Guid id,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            var result = await repository.FindStudentAsync(id, cancellationToken);
            return result is null ? Results.NotFound(new ApiError("Record not found")) : Results.Ok(result);
        });

        routes.MapPost("", async (
            StudentRequest request,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            if (string.IsNullOrWhiteSpace(request.Name) || string.IsNullOrWhiteSpace(request.Group))
            {
                return Results.BadRequest(new ApiError("Invalid students"));
            }

            var result = await repository.SaveStudentAsync(null, request, cancellationToken);
            return Results.Created($"/api/students/{result!.Id}", result);
        });

        routes.MapPut("/{id:guid}", async (
            Guid id,
            StudentRequest request,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            if (string.IsNullOrWhiteSpace(request.Name) || string.IsNullOrWhiteSpace(request.Group))
            {
                return Results.BadRequest(new ApiError("Invalid students"));
            }

            var result = await repository.SaveStudentAsync(id, request, cancellationToken);
            return result is null ? Results.NotFound(new ApiError("Record not found")) : Results.Ok(result);
        });

        routes.MapDelete("/{id:guid}", async (
            Guid id,
            CatalogRepository repository,
            CancellationToken cancellationToken) =>
        {
            try
            {
                var deleted = await repository.DeleteStudentAsync(id, cancellationToken);
                return deleted ? Results.NoContent() : Results.NotFound(new ApiError("Record not found"));
            }
            catch (PostgresException exception) when (exception.SqlState == "23503")
            {
                return Results.Conflict(new ApiError("Record is used by a lesson"));
            }
        });
    }
}
