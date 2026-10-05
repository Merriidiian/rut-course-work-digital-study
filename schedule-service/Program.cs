using Npgsql;
using University;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddSingleton<LessonRepository>();
builder.Services.AddScoped<LessonSaga>();
builder.Services.AddHttpClient().ConfigureHttpClientDefaults(b =>
    b.ConfigureHttpClient(c => c.Timeout = TimeSpan.FromSeconds(5)));
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
var app = builder.Build();
app.UseSwagger();
app.MapGet("/api/subjects", async (LessonRepository repo, CancellationToken ct) =>
{
    await using var db = await repo.OpenAsync(ct);
    await using var query = new NpgsqlCommand("SELECT id,name,hours FROM subjects ORDER BY name", db);
    await using var reader = await query.ExecuteReaderAsync(ct);
    var result = new List<Subject>();
    while (await reader.ReadAsync(ct)) result.Add(new(reader.GetGuid(0), reader.GetString(1), reader.GetInt32(2)));
    return Results.Ok(result);
});
app.MapPost("/api/subjects", async (SubjectCommand c, LessonRepository repo, CancellationToken ct) =>
{
    if (string.IsNullOrWhiteSpace(c.Name) || c.Hours <= 0) return Results.BadRequest("Invalid subject");
    await using var db = await repo.OpenAsync(ct);
    await using var insert = new NpgsqlCommand("INSERT INTO subjects(name,hours) VALUES (@name,@hours) RETURNING id", db);
    insert.Parameters.AddWithValue("name", c.Name);
    insert.Parameters.AddWithValue("hours", c.Hours);
    var id = (Guid)(await insert.ExecuteScalarAsync(ct))!;
    return Results.Created($"/api/subjects/{id}", new Subject(id, c.Name, c.Hours));
});
app.MapPut("/api/subjects/{id:guid}", async (Guid id, SubjectCommand c, LessonRepository repo, CancellationToken ct) =>
{
    if (string.IsNullOrWhiteSpace(c.Name) || c.Hours <= 0) return Results.BadRequest("Invalid subject");
    await using var db = await repo.OpenAsync(ct);
    await using var update = new NpgsqlCommand("UPDATE subjects SET name=@name,hours=@hours WHERE id=@id", db);
    update.Parameters.AddWithValue("id", id);
    update.Parameters.AddWithValue("name", c.Name);
    update.Parameters.AddWithValue("hours", c.Hours);
    return await update.ExecuteNonQueryAsync(ct) == 0 ? Results.NotFound() : Results.Ok(new Subject(id, c.Name, c.Hours));
});
app.MapDelete("/api/subjects/{id:guid}", async (Guid id, LessonRepository repo, CancellationToken ct) =>
{
    await using var db = await repo.OpenAsync(ct);
    await using var delete = new NpgsqlCommand("DELETE FROM subjects WHERE id=@id", db);
    delete.Parameters.AddWithValue("id", id);
    try { return await delete.ExecuteNonQueryAsync(ct) == 0 ? Results.NotFound() : Results.NoContent(); }
    catch (PostgresException) { return Results.Conflict("Subject is referenced"); }
});
app.MapGet("/api/lessons", async (string? group, Guid? teacherId, LessonRepository repo, CancellationToken ct) =>
    Results.Ok(await repo.GetAsync(group, teacherId, ct)));
app.MapGet("/api/students/{id:guid}/schedule", async (Guid id, IHttpClientFactory factory,
    IConfiguration config, LessonRepository repo, CancellationToken ct) =>
{
    using var response = await factory.CreateClient().GetAsync($"{config["StudentsUrl"]}/api/students/{id}", ct);
    if (!response.IsSuccessStatusCode) return Results.NotFound();
    var student = await response.Content.ReadFromJsonAsync<Student>(ct);
    return Results.Ok(await repo.GetAsync(student!.Group, null, ct));
});
app.MapPost("/api/lessons", async (LessonCommand c, LessonSaga saga, CancellationToken ct) =>
{
    try
    {
        var result = await saga.ExecuteAsync(c, ct);
        return result.Success ? Results.Created($"/api/lessons/{result.LessonId}", result) : Results.BadRequest(result);
    }
    catch (ArgumentException e) { return Results.BadRequest(e.Message); }
});
app.MapDelete("/api/lessons/{id:guid}", async (Guid id, LessonRepository repo,
    IHttpClientFactory factory, IConfiguration config, CancellationToken ct) =>
{
    using var response = await factory.CreateClient().DeleteAsync($"{config["RoomsUrl"]}/api/bookings/{id}", ct);
    if (!response.IsSuccessStatusCode) return Results.StatusCode(502);
    await using var db = await repo.OpenAsync(ct);
    await using var delete = new NpgsqlCommand("DELETE FROM lessons WHERE id=@id", db);
    delete.Parameters.AddWithValue("id", id);
    return await delete.ExecuteNonQueryAsync(ct) == 0 ? Results.NotFound() : Results.NoContent();
});
app.Run();
