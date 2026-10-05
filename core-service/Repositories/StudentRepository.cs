using Npgsql;
using University.Core.Models;

namespace University.Core.Repositories;

public partial class CatalogRepository
{
    public async Task<List<Student>> GetStudentsAsync(CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        await using var command = new NpgsqlCommand(
            "SELECT id, name, group_name, email FROM students ORDER BY name", connection);
        await using var reader = await command.ExecuteReaderAsync(cancellationToken);
        var result = new List<Student>();

        while (await reader.ReadAsync(cancellationToken))
        {
            result.Add(new Student(reader.GetGuid(0), reader.GetString(1), reader.GetString(2), reader.GetString(3)));
        }

        return result;
    }

    public async Task<Student?> FindStudentAsync(Guid id, CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        await using var command = new NpgsqlCommand(
            "SELECT id, name, group_name, email FROM students WHERE id = @id", connection);
        command.Parameters.AddWithValue("id", id);
        await using var reader = await command.ExecuteReaderAsync(cancellationToken);

        return await reader.ReadAsync(cancellationToken)
            ? new Student(reader.GetGuid(0), reader.GetString(1), reader.GetString(2), reader.GetString(3))
            : null;
    }

    public async Task<Student?> SaveStudentAsync(
        Guid? id, StudentRequest request, CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        var sql = id.HasValue
            ? "UPDATE students SET name = @name, group_name = @group_name, email = @email WHERE id = @id RETURNING id"
            : "INSERT INTO students (name, group_name, email) VALUES (@name, @group_name, @email) RETURNING id";
        await using var command = new NpgsqlCommand(sql, connection);
        command.Parameters.AddWithValue("name", request.Name);
        command.Parameters.AddWithValue("group_name", request.Group);
        command.Parameters.AddWithValue("email", request.Email);

        if (id.HasValue)
        {
            command.Parameters.AddWithValue("id", id.Value);
        }

        var result = await command.ExecuteScalarAsync(cancellationToken);
        return result is Guid savedId ? new Student(savedId, request.Name, request.Group, request.Email) : null;
    }

    public async Task<bool> DeleteStudentAsync(Guid id, CancellationToken cancellationToken)
    {
        await using var connection = await OpenAsync(cancellationToken);
        await using var command = new NpgsqlCommand("DELETE FROM students WHERE id = @id", connection);
        command.Parameters.AddWithValue("id", id);
        return await command.ExecuteNonQueryAsync(cancellationToken) > 0;
    }
}
