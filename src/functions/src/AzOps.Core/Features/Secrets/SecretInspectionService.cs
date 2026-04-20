using AzOps.Core.Common;
using System.Net;

namespace AzOps.Core.Features.Secrets;

public sealed class SecretInspectionService(ISecretValueReader secretValueReader) : ISecretInspectionService
{
    public async Task<SecretInspectionResult> InspectAsync(string secretName, CancellationToken cancellationToken)
    {
        if (string.IsNullOrWhiteSpace(secretName))
        {
            throw new AzOpsException(
                statusCode: HttpStatusCode.BadRequest,
                code: "invalid_secret_name",
                message: "A secret name is required.");
        }

        var secret = await secretValueReader.GetSecretAsync(secretName.Trim(), cancellationToken);

        return new SecretInspectionResult(
            Name: secret.Name,
            Source: "KeyVault",
            Retrieved: true,
            ValueLength: secret.Value.Length,
            ValuePreview: CreatePreview(secret.Value),
            Version: secret.Version,
            Message: "Secret retrieved successfully.");
    }

    private static string CreatePreview(string value)
    {
        if (string.IsNullOrEmpty(value))
        {
            return "(empty)";
        }

        if (value.Length <= 4)
        {
            return new string('*', value.Length);
        }

        return $"{value[..2]}...{value[^2..]}";
    }
}
