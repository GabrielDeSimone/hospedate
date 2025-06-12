import { Html, Head, Main, NextScript } from 'next/document'

export default function Document() {
    return (
        <Html className="h-full bg-white" lang="es">
            <Head>
                <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@48,400,0,0" />
                <link rel="preconnect" href="https://fonts.googleapis.com" />
                <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin />
                <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;700;800&display=swap" rel="stylesheet"/>
            </Head>
            <body className="h-full">
                <Main />
                <NextScript />
            </body>
        </Html>
    )
}