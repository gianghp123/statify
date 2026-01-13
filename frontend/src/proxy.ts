import { decodeJwt } from 'jose';
import type { NextRequest } from 'next/server';
import { NextResponse } from 'next/server';

const publicPaths = ['/', '/login', '/register', '/docs'];


export async function proxy(request: NextRequest) {
  const loginUrl = new URL('/login', request.url);
  const { pathname } = request.nextUrl;

  const isPublicPath = publicPaths.some((path) => pathname === path);

  const token = request.cookies.get('auth-token')?.value;

  if (!token && !isPublicPath) {
    return NextResponse.redirect(loginUrl);
  }

  if (token && !isPublicPath) {
    try {
      const claims = decodeJwt(token)
      if (claims.exp) {
        const isExpired = claims.exp * 1000 < Date.now()

        if (isExpired) {
          return NextResponse.redirect(loginUrl)
        }
      }
      return NextResponse.next()
    } catch (error) {
      const res = NextResponse.redirect(loginUrl);
      res.cookies.delete('auth-token');
      return res;
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
}