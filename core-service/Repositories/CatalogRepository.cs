using Npgsql;

namespace University.Core.Repositories;

public partial class CatalogRepository(IConfiguration configuration)
{
    public async Task<NpgsqlConnection> OpenAsync(CancellationToken cancellationToken)
    {
        var connection = new NpgsqlConnection(configuration.GetConnectionString("Postgres"));
        await connection.OpenAsync(cancellationToken);
        return connection;
    }
}
