namespace AzOps.Core.Features.Secrets;

public interface ISecretValueReader
{
    Task<SecretRecord> GetSecretAsync(string secretName, CancellationToken cancellationToken);
}
