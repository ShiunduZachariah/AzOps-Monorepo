using AzOps.Core.Common;
using AzOps.Infrastructure.Configuration;
using AzOps.Infrastructure.Features.Resources;
using Azure.Core;
using System.Net;

namespace AzOps.IntegrationTests.Features.Resources;

public sealed class AzureResourceGroupReaderTests
{
    [Fact]
    public async Task ListAsync_RequiresSubscriptionConfiguration()
    {
        var reader = new AzureResourceGroupReader(new FakeTokenCredential(), new AzOpsFunctionsOptions());

        var exception = await Assert.ThrowsAsync<AzOpsException>(() => reader.ListAsync(CancellationToken.None));

        Assert.Equal(HttpStatusCode.ServiceUnavailable, exception.StatusCode);
        Assert.Equal("subscription_not_configured", exception.Code);
    }

    private sealed class FakeTokenCredential : TokenCredential
    {
        public override AccessToken GetToken(TokenRequestContext requestContext, CancellationToken cancellationToken)
        {
            return new AccessToken("integration-test-token", DateTimeOffset.UtcNow.AddMinutes(5));
        }

        public override ValueTask<AccessToken> GetTokenAsync(TokenRequestContext requestContext, CancellationToken cancellationToken)
        {
            return ValueTask.FromResult(GetToken(requestContext, cancellationToken));
        }
    }
}
