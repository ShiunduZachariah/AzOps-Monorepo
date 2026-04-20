namespace AzOps.Core.Common;

public sealed record ErrorResponse(string Code, string Message, string? Details = null);
