export const metadata = {
  title: 'Next.js Bun Test',
  description: 'Testing Next.js with Bun package manager',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
