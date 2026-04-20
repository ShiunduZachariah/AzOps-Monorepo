namespace AzOps.Infrastructure.Configuration;

public sealed class AzOpsFunctionsOptions
{
    public string? SubscriptionId { get; init; }

    public string? KeyVaultName { get; init; }

    public string? KeyVaultUri { get; init; }

    public Uri? ResolveKeyVaultUri()
    {
        var configuredUri = KeyVaultUri?.Trim();
        if (!string.IsNullOrWhiteSpace(configuredUri))
        {
            return new Uri(configuredUri, UriKind.Absolute);
        }

        var configuredName = KeyVaultName?.Trim();
        if (!string.IsNullOrWhiteSpace(configuredName))
        {
            return new Uri($"https://{configuredName}.vault.azure.net/", UriKind.Absolute);
        }

        return null;
    }
}
