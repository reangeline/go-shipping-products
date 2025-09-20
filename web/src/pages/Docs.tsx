export default function Docs() {
  // Em dev, proxie /docs para o backend no vite.config.ts
  // server.proxy = { '/v1': 'http://localhost:8080', '/docs': 'http://localhost:8080' }
  return (
    <div className="h-[75vh]">
      <iframe
        title="Swagger UI"
        src="/docs"
        className="w-full h-full rounded border"
      />
    </div>
  );
}