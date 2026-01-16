'use client'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { UserDto } from "@/features/users/dtos/response/user.response.dto"

interface ProfileSettingCardProps {
  user: UserDto
}

export function ProfileSettingCard({ user }: ProfileSettingCardProps) {
  return <Card className="glass-panel border-border/40">
    <CardHeader>
      <CardTitle>Public Profile</CardTitle>
      <CardDescription>
        How you appear to other members.
      </CardDescription>
    </CardHeader>
    <CardContent className="flex gap-6">
      <div className="flex items-center gap-6">
        <Avatar className="size-20 border-2 border-primary/20">
          <AvatarImage src="" />
          <AvatarFallback className="bg-primary/10 text-primary text-xl font-bold">{user.username[0].toUpperCase()}</AvatarFallback>
        </Avatar>
      </div>

      <div className="flex flex-col gap-4 flex-1">
        <div className="space-y-2 flex-1">
          <Label htmlFor="username">Username</Label>
          <Input id="username" defaultValue={`${user.username}`} className="bg-background/50 border-border/40" />
        </div>
        <div className="space-y-2 flex-1">
          <Label htmlFor="email">Email Address</Label>
          <Input id="email" type="email" defaultValue={`${user.email}`} className="bg-background/50 border-border/40" />
        </div>
        {/* <div className="space-y-2">
          <Label htmlFor="fullname">Full Name</Label>
          <Input id="fullname" defaultValue={`${user.fullname}`} className="bg-background/50 border-border/40" />
        </div>
        <div className="space-y-2">
          <Label htmlFor="location">Location</Label>
          <Input id="location" placeholder="San Francisco, CA" defaultValue={`${user.location}`} className="bg-background/50 border-border/40" />
        </div> */}
      </div>
    </CardContent>
    <CardFooter className="justify-end border-t border-border/40 pt-6">
      <Button className="shadow-neon">Update Profile</Button>
    </CardFooter>
  </Card>
} 