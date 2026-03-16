export const INITIAL_SPOTS = Array.from({ length: 20 }, (_, i) => {
  const row = Math.floor(i / 10);
  const col = i % 10;
  return {
    id: `spot-${i + 1}`,
    label: `A-${i + 1}`,
    x: 50 + col * 70,
    y: 50 + row * 150,
    width: 60,
    height: 100,
  };
});
