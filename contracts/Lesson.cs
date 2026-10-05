namespace University.Contracts;

public record Lesson(
    Guid Id,
    Guid TeacherId,
    Guid SubjectId,
    Guid RoomId,
    string Group,
    DateTimeOffset StartsAt,
    DateTimeOffset EndsAt,
    string Status);
