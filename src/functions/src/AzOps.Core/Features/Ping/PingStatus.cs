namespace AzOps.Core.Features.Ping;

public sealed record PingStatus(string Status, string Service, string Message, DateTimeOffset TimestampUtc);
