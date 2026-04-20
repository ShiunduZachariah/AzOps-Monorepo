using AzOps.Core.Common;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Extensions.Logging;
using System.Net;

namespace AzOps.Functions.Features.Shared;

internal static class FunctionExecution
{
    public static async Task<HttpResponseData> ExecuteAsync<T>(
        HttpRequestData request,
        ILogger logger,
        string operationName,
        Func<CancellationToken, Task<T>> handler,
        CancellationToken cancellationToken = default)
    {
        try
        {
            var payload = await handler(cancellationToken);
            return await CreateJsonResponseAsync(request, HttpStatusCode.OK, payload, cancellationToken);
        }
        catch (AzOpsException ex)
        {
            logger.LogWarning(ex, "{OperationName} failed with code {Code}", operationName, ex.Code);
            return await CreateJsonResponseAsync(
                request,
                ex.StatusCode,
                new ErrorResponse(ex.Code, ex.Message, ex.Details),
                cancellationToken);
        }
        catch (Exception ex)
        {
            logger.LogError(ex, "{OperationName} failed unexpectedly", operationName);
            return await CreateJsonResponseAsync(
                request,
                HttpStatusCode.InternalServerError,
                new ErrorResponse("unexpected_error", "An unexpected error occurred."),
                cancellationToken);
        }
    }

    private static async Task<HttpResponseData> CreateJsonResponseAsync<T>(
        HttpRequestData request,
        HttpStatusCode statusCode,
        T payload,
        CancellationToken cancellationToken)
    {
        var response = request.CreateResponse(statusCode);
        await response.WriteAsJsonAsync(payload, cancellationToken);
        return response;
    }
}
