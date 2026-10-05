using University.Core.Endpoints;
using University.Core.Repositories;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddSingleton<CatalogRepository>();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();
app.UseSwagger();

app.MapStudentEndpoints();

app.MapTeacherEndpoints();

app.MapSubjectEndpoints();
app.Run();
