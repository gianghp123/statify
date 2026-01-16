import { logout } from "@/features/auth/services/auth.actions";
import { getCurrentUser } from "@/features/auth/services/auth.get";
import { PreferenceSettingCard } from "@/features/settings/components/PreferenceSettingCard";
import { ProfileSettingCard } from "@/features/settings/components/ProfileSettingCard";

export default async function AccountSettings() {
  const userRes = await getCurrentUser()

  if (!userRes.success || !userRes.data) {
    logout()
    return
  }
  const user = userRes.data

  return (
    <div className="max-w-4xl space-y-8">
      <div>
        <h1 className="text-foreground tracking-tight text-3xl font-bold">Account Settings</h1>
        <p className="text-muted-foreground mt-1">
          Manage your personal information and account preferences.
        </p>
      </div>
      <ProfileSettingCard user={user} />
      <PreferenceSettingCard />
    </div>
  );
}
