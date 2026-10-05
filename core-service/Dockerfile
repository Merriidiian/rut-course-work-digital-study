FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src
COPY contracts/University.Contracts.csproj ./contracts/
COPY core-service/Core.csproj core-service/NuGet.Config ./core-service/
RUN dotnet restore core-service/Core.csproj --configfile core-service/NuGet.Config -p:NuGetAudit=false
COPY contracts ./contracts
COPY core-service ./core-service
RUN dotnet publish core-service/Core.csproj -c Release -o /app --no-restore
FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=build /app .
ENV ASPNETCORE_URLS=http://+:8080
ENTRYPOINT ["dotnet", "Core.dll"]
