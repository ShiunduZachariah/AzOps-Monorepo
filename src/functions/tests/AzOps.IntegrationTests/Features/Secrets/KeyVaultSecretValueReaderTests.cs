using AzOps.Core.Common;
using AzOps.Infrastructure.Configuration;
using AzOps.Infrastructure.Features.Secrets;
using Azure.Core;
using System.Net;

namespace AzOps.IntegrationTests.Features.Secrets;

public sealed class KeyVaultSecretValueReaderTests
{
    [Fact]
    public async Task GetSecretAsync_RequiresKeyVaultConfiguration()
    {
        var reader = new KeyVaultSecretValueReader(new FakeTokenCredential(), new AzOpsFunctionsOptions());

        var exception = await Assert.ThrowsAsync<AzOpsException>(() => reader.GetSecretAsync("DbPassword", CancellationToken.None));

        Assert.Equal(HttpStatusCode.ServiceUnavailable, exception.StatusCode);
        Assert.Equal("key_vault_not_configured", exception.Code);
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
