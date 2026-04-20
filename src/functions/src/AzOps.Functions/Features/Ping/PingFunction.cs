using AzOps.Core.Features.Ping;
using AzOps.Functions.Features.Shared;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Extensions.Logging;

namespace AzOps.Functions.Features.Ping;

public sealed class PingFunction(IPingService pingService, ILogger<PingFunction> logger)
{
    [Function("Ping")]
    public Task<HttpResponseData> Run(
        [HttpTrigger(AuthorizationLevel.Anonymous, "get", Route = "ping")] HttpRequestData request,
        FunctionContext context)
    {
        return FunctionExecution.ExecuteAsync(
            request,
            logger,
            operationName: "ping",
            handler: _ => Task.FromResult(pingService.GetStatus()),
            cancellationToken: context.CancellationToken);
    }
}
