using Mainframe.Server.Auth;
using Mainframe.Server.Recipes;

using Scalar.AspNetCore;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddOpenApi();
builder.Services.AddAuthModule();
builder.Services.AddRecipesModule();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
    app.MapScalarApiReference();
    app.UseWebAssemblyDebugging();
}

app.UseHttpsRedirection();
app.UseStaticFiles();
app.UseBlazorFrameworkFiles();

app.MapAuthEndpoints();
app.MapRecipesEndpoints();

app.MapFallbackToFile("index.html");
app.Run();
