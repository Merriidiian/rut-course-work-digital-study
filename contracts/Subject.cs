namespace University.Contracts;

public record Subject(Guid Id, string Name, int Hours);

public record SubjectRequest(string Name, int Hours);
