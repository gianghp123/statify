export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center relative overflow-hidden selection:bg-primary selection:text-black text-white px-4">
      {/* Background Glows */}
      <div
        className="absolute top-[-20%] left-[-10%] w-[600px] h-[600px] bg-purple-900/40 rounded-full blur-[120px] pointer-events-none"
        aria-hidden="true"
      ></div>
      <div
        className="absolute bottom-[-20%] right-[-10%] w-[600px] h-[600px] bg-primary/10 rounded-full blur-[120px] pointer-events-none"
        aria-hidden="true"
      ></div>
      
      <main className="w-full max-w-[380px] relative z-10">
        {children}
        <div className="text-center mt-8">
          <p className="text-[#bfacc8]/30 text-xs text-center">
            Statify Deep Plum Edition © 2024
          </p>
        </div>
      </main>
    </div>
  );
}
