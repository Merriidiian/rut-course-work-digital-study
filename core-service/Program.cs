using University.Core.Endpoints;
using University.Core.Clients;
using University.Core.Repositories;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddSingleton<CatalogRepository>();
builder.Services.AddHttpClient<ScheduleClient>(client =>
{
    client.BaseAddress = new Uri(builder.Configuration["ScheduleUrl"]!);
    client.Timeout = TimeSpan.FromSeconds(5);
});
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();
app.UseSwagger();

app.MapStudentEndpoints();

app.MapTeacherEndpoints();

app.MapSubjectEndpoints();
app.Run();
