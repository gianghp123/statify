import { Settings as SettingsIcon, User, Shield, CreditCard, Bell, Globe } from "lucide-react";

export default function SettingsPage() {
  const sections = [
    { title: "General", icon: SettingsIcon, description: "Manage your project settings and configuration." },
    { title: "Account", icon: User, description: "Update your personal information and preferences." },
    { title: "Security", icon: Shield, description: "Configure two-factor authentication and passwords." },
    { title: "Billing", icon: CreditCard, description: "Manage your subscription, invoices, and payments." },
    { title: "Notifications", icon: Bell, description: "Control how you receive updates and alerts." },
    { title: "Domains", icon: Globe, description: "Configure custom domains and SSL certificates." },
  ];

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-foreground tracking-tight text-3xl md:text-4xl font-bold leading-tight">
          Settings
        </h1>
        <p className="text-muted-foreground mt-1 text-sm md:text-base">
          Manage your account and site preferences.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {sections.map((section) => (
          <div key={section.title} className="bg-card rounded-xl p-6 border border-border hover:border-primary/30 transition-all group cursor-pointer shadow-lg">
            <div className="size-12 rounded-lg bg-primary/10 flex items-center justify-center text-primary mb-4 group-hover:scale-110 transition-transform">
              <section.icon className="w-6 h-6" />
            </div>
            <h3 className="text-foreground font-bold text-lg mb-2">{section.title}</h3>
            <p className="text-muted-foreground text-sm leading-relaxed">
              {section.description}
            </p>
          </div>
        ))}
      </div>

      <footer className="mt-12 mb-6 flex justify-center text-xs text-muted-foreground/50">
        <p>© 2024 Statify Inc. Custom Deep Plum Edition.</p>
      </footer>
    </div>
  );
}
