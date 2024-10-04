import { ActionIcon, Card, Flex, Group, Image, Input, Menu, Title, useComputedColorScheme, useMantineColorScheme } from '@mantine/core';
import { Cog, EllipsisVertical, LogOut, Moon, Search, Sun } from "lucide-react";
import { useMeta } from '../state/use-meta';
import { useAuth } from '../state/use-auth';
import useFetch from '../hooks/useFetch';

export default function Header() {
  const { toggleColorScheme } = useMantineColorScheme();
  const isDark = useComputedColorScheme() === "dark";
  const { pageTitle } = useMeta();
  const { loggedIn } = useAuth();
  const fetchData = useFetch();

  async function handleLogout() {
    await fetchData("/api/mg/logout");
    window.location.reload();
  }

  return (
    <Card
      withBorder
      radius={0}
      shadow='md'
      style={{ position: 'sticky', top: 0, zIndex: 100 }}
    >
      <Flex
        p={0}
        justify="space-between"
      >
        <Group gap="xs">
          <Image
            w={{ xs: 36, sm: 54 }}
            src="/media/d2m.png"
            alt="Logo of D2M"
          />
          <Title
            order={2}
            fz={{ xs: 18, sm: 24 }}
          >
            {pageTitle}
          </Title>
        </Group>

        <Group gap="xs">
          {loggedIn && (
            <Input
              readOnly
              radius={50}
              variant='filled'
              placeholder='Search...'
              w={{ xs: 115, sm: 180 }}
              leftSection={<Search />}
            />
          )}
          <Menu shadow="md" width={200}>
            <Menu.Target>
              <ActionIcon
                aria-label='Menu button'
                variant='light'
                size="lg"
                py={18}
              >
                <EllipsisVertical />
              </ActionIcon>
            </Menu.Target>

            <Menu.Dropdown>
              <Menu.Label>Application</Menu.Label>
              <Menu.Item
                leftSection={isDark ?  <Sun /> : <Moon />}
                onClick={toggleColorScheme}
              >
                Switch to {isDark ? "light" : "dark"} theme
              </Menu.Item>
              {loggedIn && (
                <>
                  <Menu.Item leftSection={<Cog />}>
                    Settings
                  </Menu.Item>
                  <Menu.Item
                    color='red'
                    onClick={handleLogout}
                    leftSection={<LogOut />}
                  >
                    Logout
                  </Menu.Item>
                </>
              )}
            </Menu.Dropdown>
          </Menu>
        </Group>
      </Flex>
    </Card>
  )
}
