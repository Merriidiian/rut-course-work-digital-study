namespace University;

public record Subject(Guid Id, string Name, int Hours);
public record SubjectCommand(string Name, int Hours);
public record LessonCommand(Guid TeacherId, Guid SubjectId, Guid RoomId,
    string Group, DateTimeOffset StartsAt, DateTimeOffset EndsAt, bool FailAfterBooking = false);
public record Lesson(Guid Id, Guid TeacherId, Guid SubjectId, Guid RoomId,
    string Group, DateTimeOffset StartsAt, DateTimeOffset EndsAt, string Status);
public record Student(Guid Id, string Name, string Group, string Email, int Capacity);
public record Teacher(Guid Id, string Name, string Group, string Email, int Capacity);
public record SagaResult(bool Success, Guid LessonId, string Message, bool CompensationExecuted);
