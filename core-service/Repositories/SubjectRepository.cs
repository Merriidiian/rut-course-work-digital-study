using Npgsql;
using University.Contracts;

namespace University.Core.Repositories;

public partial class CatalogRepository
{
    public async Task<List<Subject>> GetSubjectsAsync(CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        await using var command = new NpgsqlCommand(
            "SELECT id, name, hours FROM subjects ORDER BY name", connection);
        await using var reader = await command.ExecuteReaderAsync(cancellationToken);
        var result = new List<Subject>();

        while (await reader.ReadAsync(cancellationToken))
        {
            result.Add(new Subject(reader.GetGuid(0), reader.GetString(1), reader.GetInt32(2)));
        }

        return result;
    }

    public async Task<Subject?> FindSubjectAsync(Guid id, CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        await using var command = new NpgsqlCommand(
            "SELECT id, name, hours FROM subjects WHERE id = @id", connection);
        command.Parameters.AddWithValue("id", id);
        await using var reader = await command.ExecuteReaderAsync(cancellationToken);

        return await reader.ReadAsync(cancellationToken)
            ? new Subject(reader.GetGuid(0), reader.GetString(1), reader.GetInt32(2))
            : null;
    }

    public async Task<Subject?> SaveSubjectAsync(
        Guid? id, SubjectRequest request, CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        var sql = id.HasValue
            ? "UPDATE subjects SET name = @name, hours = @hours WHERE id = @id RETURNING id"
            : "INSERT INTO subjects (name, hours) VALUES (@name, @hours) RETURNING id";
        await using var command = new NpgsqlCommand(sql, connection);
        command.Parameters.AddWithValue("name", request.Name);
        command.Parameters.AddWithValue("hours", request.Hours);

        if (id.HasValue)
        {
            command.Parameters.AddWithValue("id", id.Value);
        }

        var result = await command.ExecuteScalarAsync(cancellationToken);
        return result is Guid savedId ? new Subject(savedId, request.Name, request.Hours) : null;
    }

    public async Task<bool> DeleteSubjectAsync(Guid id, CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        await using var command = new NpgsqlCommand("DELETE FROM subjects WHERE id = @id", connection);
        command.Parameters.AddWithValue("id", id);
        return await command.ExecuteNonQueryAsync(cancellationToken) > 0;
    }
}
