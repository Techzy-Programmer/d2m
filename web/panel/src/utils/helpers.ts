export function formatDate(date: string | number) {
  return (new Date(date).toLocaleString('en-GB', {
    weekday: 'short',
    year: '2-digit',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    timeZone: "IST",
    hour12: true
  })) + " IST";
}
