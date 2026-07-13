const capabilities = [
  "Repository analysis",
  "Code quality insights",
  "Security scanning",
  "AI documentation",
];

export function App() {
  return (
    <main className="min-h-screen bg-zinc-950 text-zinc-50">
      <section className="mx-auto flex min-h-screen w-full max-w-6xl flex-col justify-center px-6 py-10">
        <p className="text-sm font-semibold uppercase tracking-wide text-emerald-300">
          Mintok
        </p>
        <h1 className="mt-4 max-w-3xl text-5xl font-semibold leading-tight">
          AI-powered developer productivity for healthier repositories.
        </h1>
        <p className="mt-6 max-w-2xl text-lg leading-8 text-zinc-300">
          Analyze codebases, surface quality and security issues, generate
          documentation, and prepare pull request insights from one workspace.
        </p>
        <div className="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {capabilities.map((capability) => (
            <div
              className="rounded-lg border border-zinc-800 bg-zinc-900 p-4 text-sm font-medium text-zinc-200"
              key={capability}
            >
              {capability}
            </div>
          ))}
        </div>
      </section>
    </main>
  );
}
