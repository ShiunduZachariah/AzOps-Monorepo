using AzOps.Core.Common;
using AzOps.Core.Features.Secrets;
using System.Net;

namespace AzOps.Functions.Tests.Features.Secrets;

public sealed class SecretInspectionServiceTests
{
    [Fact]
    public async Task InspectAsync_ReturnsMaskedSecretPreview()
    {
        var service = new SecretInspectionService(new FakeSecretValueReader("DbPassword", "supersafe", "v1"));

        var result = await service.InspectAsync("DbPassword", CancellationToken.None);

        Assert.Equal("DbPassword", result.Name);
        Assert.True(result.Retrieved);
        Assert.Equal(9, result.ValueLength);
        Assert.Equal("su...fe", result.ValuePreview);
        Assert.Equal("v1", result.Version);
    }

    [Fact]
    public async Task InspectAsync_RejectsBlankSecretNames()
    {
        var service = new SecretInspectionService(new FakeSecretValueReader("ignored", "ignored", null));

        var exception = await Assert.ThrowsAsync<AzOpsException>(() => service.InspectAsync("   ", CancellationToken.None));

        Assert.Equal(HttpStatusCode.BadRequest, exception.StatusCode);
        Assert.Equal("invalid_secret_name", exception.Code);
    }

    private sealed class FakeSecretValueReader(string name, string value, string? version) : ISecretValueReader
    {
        public Task<SecretRecord> GetSecretAsync(string secretName, CancellationToken cancellationToken)
        {
            return Task.FromResult(new SecretRecord(name, value, version));
        }
    }
}
