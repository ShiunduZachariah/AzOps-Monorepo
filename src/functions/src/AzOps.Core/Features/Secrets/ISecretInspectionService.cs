namespace AzOps.Core.Features.Secrets;

public interface ISecretInspectionService
{
    Task<SecretInspectionResult> InspectAsync(string secretName, CancellationToken cancellationToken);
}
