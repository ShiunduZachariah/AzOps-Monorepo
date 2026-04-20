namespace AzOps.Core.Features.Secrets;

public sealed record SecretRecord(string Name, string Value, string? Version);
