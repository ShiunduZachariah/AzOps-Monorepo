namespace AzOps.Core.Features.Secrets;

public sealed record SecretInspectionResult(
    string Name,
    string Source,
    bool Retrieved,
    int ValueLength,
    string ValuePreview,
    string? Version,
    string Message);
