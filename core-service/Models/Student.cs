namespace University.Core.Models;

public record Student(Guid Id, string Name, string Group, string Email);

public record StudentRequest(string Name, string Group, string Email);
