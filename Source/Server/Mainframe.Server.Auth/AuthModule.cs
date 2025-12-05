using System.Reflection;

using Mainframe.Server.Auth.Features.Passwords;
using Mainframe.Server.Infrastructure;

using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Routing;
using Microsoft.Extensions.DependencyInjection;

namespace Mainframe.Server.Auth;

public static class AuthModule
{
    public static IServiceCollection AddAuthModule(this IServiceCollection services)
    {
        services.AddScoped<IPasswordHasher, Argon2IdPasswordHasher>();
        return services;
    }

    public static IEndpointRouteBuilder MapAuthEndpoints(this IEndpointRouteBuilder app)
    {
        var group = app.MapGroup("/api/auth")
            .WithTags("Auth");

        group.MapEndpointsFromAssembly(Assembly.GetExecutingAssembly());

        return app;
    }
}