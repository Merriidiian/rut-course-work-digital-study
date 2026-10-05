using University.Contracts;

namespace University.Core.Clients;

public class ScheduleClient(HttpClient httpClient)
{
    public Task<bool> HasTeacherLessonsAsync(Guid id, CancellationToken cancellationToken)
    {
        return HasLessonsAsync($"api/lessons?teacherId={id}", cancellationToken);
    }

    public Task<bool> HasSubjectLessonsAsync(Guid id, CancellationToken cancellationToken)
    {
        return HasLessonsAsync($"api/lessons?subjectId={id}", cancellationToken);
    }

    private async Task<bool> HasLessonsAsync(string path, CancellationToken cancellationToken)
    {
        var lessons = await httpClient.GetFromJsonAsync<List<Lesson>>(path, cancellationToken);
        return lessons is not null && lessons.Count > 0;
    }
}
