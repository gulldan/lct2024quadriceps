import "./globals.css";
import "bootstrap/dist/css/bootstrap.min.css";
import { Links } from "./ui/links";
//const inter = Inter({ subsets: ["latin"] });

export const metadata = {
  title: "КОПИРАЙТРШТНАЯ",
  description: "Мы все видим",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <head>
        <link rel="manifest" href="/manifest.json" />
        <link rel="icon" href="/gool2.png" sizes="any" />
        <link rel="apple-touch-icon" href="/gool2.png"/>
      </head>
      <body>
        <Links />
        {children}
      </body>
    </html>
  );
}

//<body className={inter.className}>
