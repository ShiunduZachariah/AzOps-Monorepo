using AzOps.Core.Common;
using AzOps.Core.Features.Secrets;
using AzOps.Infrastructure.Configuration;
using Azure;
using Azure.Core;
using Azure.Security.KeyVault.Secrets;
using System.Net;

namespace AzOps.Infrastructure.Features.Secrets;

public sealed class KeyVaultSecretValueReader : ISecretValueReader
{
    private readonly TokenCredential credential;
    private readonly AzOpsFunctionsOptions options;
    private SecretClient? secretClient;

    public KeyVaultSecretValueReader(TokenCredential credential, AzOpsFunctionsOptions options)
    {
        this.credential = credential;
        this.options = options;
    }

    public async Task<SecretRecord> GetSecretAsync(string secretName, CancellationToken cancellationToken)
    {
        try
        {
            var client = GetSecretClient();
            Response<KeyVaultSecret> response = await client.GetSecretAsync(secretName, cancellationToken: cancellationToken);
            KeyVaultSecret secret = response.Value;

            return new SecretRecord(
                Name: secret.Name,
                Value: secret.Value,
                Version: secret.Properties.Version);
        }
        catch (RequestFailedException ex) when (ex.Status == 404)
        {
            throw new AzOpsException(
                statusCode: HttpStatusCode.NotFound,
                code: "secret_not_found",
                message: $"Secret '{secretName}' was not found in Key Vault.",
                details: ex.Message,
                innerException: ex);
        }
        catch (RequestFailedException ex)
        {
            throw new AzOpsException(
                statusCode: HttpStatusCode.BadGateway,
                code: "key_vault_request_failed",
                message: "The Key Vault request failed.",
                details: ex.Message,
                innerException: ex);
        }
    }

    private SecretClient GetSecretClient()
    {
        if (secretClient is not null)
        {
            return secretClient;
        }

        var vaultUri = options.ResolveKeyVaultUri();
        if (vaultUri is null)
        {
            throw new AzOpsException(
                statusCode: HttpStatusCode.ServiceUnavailable,
                code: "key_vault_not_configured",
                message: "AZOPS_KEY_VAULT_URI or AZOPS_KEY_VAULT_NAME must be configured for secret retrieval.");
        }

        secretClient = new SecretClient(vaultUri, credential);
        return secretClient;
    }
}
