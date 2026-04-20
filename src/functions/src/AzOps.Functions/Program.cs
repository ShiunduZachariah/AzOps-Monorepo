using AzOps.Infrastructure.DependencyInjection;
using AzOps.Infrastructure.Configuration;
using Microsoft.Azure.Functions.Worker.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using System.Text.Json;

var builder = FunctionsApplication.CreateBuilder(args);

builder.ConfigureFunctionsWebApplication();

builder.Services.Configure<JsonSerializerOptions>(options =>
{
    options.PropertyNamingPolicy = JsonNamingPolicy.CamelCase;
});

var azOpsOptions = new AzOpsFunctionsOptions
{
    SubscriptionId = builder.Configuration["AZOPS_SUBSCRIPTION_ID"] ?? builder.Configuration["AZURE_SUBSCRIPTION_ID"],
    KeyVaultName = builder.Configuration["AZOPS_KEY_VAULT_NAME"],
    KeyVaultUri = builder.Configuration["AZOPS_KEY_VAULT_URI"]
};

builder.Services.AddAzOpsInfrastructure(azOpsOptions);

var app = builder.Build();
app.Run();
