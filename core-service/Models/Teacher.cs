namespace University.Core.Models;

public record Teacher(Guid Id, string Name, string Email);

public record TeacherRequest(string Name, string Email);
