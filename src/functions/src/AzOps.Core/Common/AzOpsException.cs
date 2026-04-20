using System.Net;

namespace AzOps.Core.Common;

public sealed class AzOpsException : Exception
{
    public AzOpsException(HttpStatusCode statusCode, string code, string message, string? details = null, Exception? innerException = null)
        : base(message, innerException)
    {
        StatusCode = statusCode;
        Code = code;
        Details = details;
    }

    public HttpStatusCode StatusCode { get; }

    public string Code { get; }

    public string? Details { get; }
}
