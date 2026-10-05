using Npgsql;
using University.Contracts;

namespace University.Core.Repositories;

public partial class CatalogRepository
{
    public async Task<List<Teacher>> GetTeachersAsync(CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        await using var command = new NpgsqlCommand(
            "SELECT id, name, email FROM teachers ORDER BY name", connection);
        await using var reader = await command.ExecuteReaderAsync(cancellationToken);
        var result = new List<Teacher>();

        while (await reader.ReadAsync(cancellationToken))
        {
            result.Add(new Teacher(reader.GetGuid(0), reader.GetString(1), reader.GetString(2)));
        }

        return result;
    }

    public async Task<Teacher?> FindTeacherAsync(Guid id, CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        await using var command = new NpgsqlCommand(
            "SELECT id, name, email FROM teachers WHERE id = @id", connection);
        command.Parameters.AddWithValue("id", id);
        await using var reader = await command.ExecuteReaderAsync(cancellationToken);

        return await reader.ReadAsync(cancellationToken)
            ? new Teacher(reader.GetGuid(0), reader.GetString(1), reader.GetString(2))
            : null;
    }

    public async Task<Teacher?> SaveTeacherAsync(
        Guid? id, TeacherRequest request, CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        var sql = id.HasValue
            ? "UPDATE teachers SET name = @name, email = @email WHERE id = @id RETURNING id"
            : "INSERT INTO teachers (name, email) VALUES (@name, @email) RETURNING id";
        await using var command = new NpgsqlCommand(sql, connection);
        command.Parameters.AddWithValue("name", request.Name);
        command.Parameters.AddWithValue("email", request.Email);

        if (id.HasValue)
        {
            command.Parameters.AddWithValue("id", id.Value);
        }

        var result = await command.ExecuteScalarAsync(cancellationToken);
        return result is Guid savedId ? new Teacher(savedId, request.Name, request.Email) : null;
    }

    public async Task<bool> DeleteTeacherAsync(Guid id, CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        await using var command = new NpgsqlCommand("DELETE FROM teachers WHERE id = @id", connection);
        command.Parameters.AddWithValue("id", id);
        return await command.ExecuteNonQueryAsync(cancellationToken) > 0;
    }
}
