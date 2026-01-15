import { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import { decodeJwt } from "jose";

const publicPaths = ['/', '/login', '/register'];
const publicPrefixes = ['/documentation'];

export async function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const loginUrl = new URL('/login', request.url);
  const token = request.cookies.get('auth-token')?.value;

  const isPublicPath = publicPaths.includes(pathname) ||
    publicPrefixes.some(prefix => pathname.startsWith(prefix));

  if (isPublicPath) {
    return NextResponse.next();
  }

  if (!token) {
    return NextResponse.redirect(loginUrl);
  }
  try {
    const claims = decodeJwt(token);
    if (claims.exp && claims.exp * 1000 < Date.now()) {
      return NextResponse.redirect(loginUrl);
    }
    return NextResponse.next();
  } catch (error) {
    const res = NextResponse.redirect(loginUrl);
    res.cookies.delete('auth-token');
    return res;
  }
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
}