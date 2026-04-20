using AzOps.Core.Features.Health;
using AzOps.Functions.Features.Shared;
using Microsoft.Azure.Functions.Worker;
using Microsoft.Azure.Functions.Worker.Http;
using Microsoft.Extensions.Logging;

namespace AzOps.Functions.Features.Health;

public sealed class HealthFunction(IHealthService healthService, ILogger<HealthFunction> logger)
{
    [Function("Health")]
    public Task<HttpResponseData> Run(
        [HttpTrigger(AuthorizationLevel.Anonymous, "get", Route = "health")] HttpRequestData request,
        FunctionContext context)
    {
        return FunctionExecution.ExecuteAsync(
            request,
            logger,
            operationName: "health",
            handler: _ => Task.FromResult(healthService.GetStatus()),
            cancellationToken: context.CancellationToken);
    }
}
