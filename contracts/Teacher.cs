namespace University.Contracts;

public record Teacher(Guid Id, string Name, string Email);

public record TeacherRequest(string Name, string Email);
