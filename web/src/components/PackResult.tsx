type Props = {
  itemsByPack: Record<number, number>;
  totalItems: number;
  totalPacks: number;
  leftover: number;
};

export default function PackResult({ itemsByPack, totalItems, totalPacks, leftover }: Props) {
  const entries = Object.entries(itemsByPack)
    .map(([size, count]) => [Number(size), count as number] as const)
    .sort((a, b) => b[0] - a[0]); // maior → menor

  return (
    <div className="rounded-lg border p-4 space-y-3">
      <h2 className="font-semibold">Resultado</h2>

      <div className="grid grid-cols-2 gap-2 text-sm">
        <div className="text-gray-600">Total de itens</div>
        <div className="font-medium">{totalItems}</div>

        <div className="text-gray-600">Total de pacotes</div>
        <div className="font-medium">{totalPacks}</div>

        <div className="text-gray-600">Sobra</div>
        <div className="font-medium">{leftover}</div>
      </div>

      <div>
        <h3 className="text-sm font-semibold mb-1">Por tamanho:</h3>
        {entries.length === 0 ? (
          <p className="text-sm text-gray-500">Nenhum pacote.</p>
        ) : (
          <ul className="text-sm list-disc pl-5">
            {entries.map(([size, count]) => (
              <li key={size}>
                {count} × {size}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}