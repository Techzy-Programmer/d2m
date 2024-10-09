import { DeploymentDetailResp } from "../utils/types";
import { useDisclosure } from "@mantine/hooks";
import { useEffect, useState } from "react";
import Loading from "./Loading";

import {
  Box,
  Text,
  Card,
  Flex,
  Group,
  Paper,
  Title,
  Space,
  Anchor,
  Center,
  Button,
  Tooltip,
  Timeline,
  Collapse,
  HoverCard,
} from "@mantine/core";

import {
  Bug,
  Clock,
  Github,
  Bookmark,
  GitBranch,
  GitCommit,
  CheckCheck,
  BadgeXIcon,
  MessageCircle,
  TriangleAlert,
  AlarmClockPlus,
  BadgeCheckIcon,
} from "lucide-react";
import { formatDate } from "../utils/helpers";

export default function Deployment() {
  const [deploymentData, setDeploymentData] = useState<DeploymentDetailResp>();
  const [opened, { toggle }] = useDisclosure(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // ToDo: Fetch actual deployment data from API

    setLoading(false);
    setDeploymentData({
      meta: {
        CommitHash: "be1457de93b09d2bc999141d96710230adb9e5e9",
        CommitMsg: "Add basic deployment details ui page",
        StartAt: "2024-10-06 19:54:15.0509262+05:30",
        EndAt: "2024-10-06 19:54:16.4255948+05:30",
        Repo: "Techzy-Programmer/d2m",
        Status: "failed",
        Branch: "main",
        ID: 90,
      },
      message: "Deployment details fetched successfully",
      ok: true,
      logs: [
        {
          ID: 1,
          Level: 0,
          DeployID: 90,
          Timestamp: 1728224655,
          Message: "Pre-deployment script execution while deploying",
          Title: "Pre-Deployment Script"
        },
        {
          ID: 1,
          Level: 1,
          DeployID: 90,
          Timestamp: 1728224655,
          Message: "Runtime data expired",
          Title: "Data Check"
        },
        {
          ID: 2,
          Level: 3,
          DeployID: 90,
          Timestamp: 1728224656,
          Message: "Deployment Process Executed",
          Title: "Deployment Exec"
        },
        {
          ID: 3,
          Level: 2,
          DeployID: 90,
          Timestamp: 1728224656,
          Message: "Deployment failed",
          Title: "Deployment Failed"
        }
      ]
    });
  }, []);

  if (loading) {
    return <Loading />;
  }

  const { meta } = deploymentData!;
  
  return (
    <>
      <Center>
        <Paper w="90%" shadow="sm" my={20} p={20}>
          <Flex gap={20} justify="space-between">
            <Box>
              <Anchor href={`http://github.com/${meta.Repo}`} target="_blank">
                <Group wrap="nowrap" gap={8} mb={10}>
                  <Github />
                  <Title style={{ wordBreak: "break-all" }} order={4}>
                    {meta.Repo}
                  </Title>
                </Group>
              </Anchor>
              <Tooltip withArrow label="click here" zIndex={10} opened position="right" offset={10} px={10} py={5}>
                <Group
                  onClick={toggle}
                  w="fit-content"
                  wrap="nowrap"
                  c="cyan"
                  gap={8}

                  style={{
                    textDecorationLine: 'underline',
                    textDecorationStyle: 'wavy',
                    textUnderlineOffset: 5,
                    wordBreak: "break-all"
                  }}
                >
                  <GitCommit />
                  <Text style={{ cursor: "pointer" }}>{meta.CommitHash.substring(0, 7)}</Text>
                </Group>
              </Tooltip>
            </Box>
            <Box>
              <HoverCard withArrow>
                <HoverCard.Target>
                  {meta.Status === "success" ? <BadgeCheckIcon color="#1ebe00" /> : <BadgeXIcon color="red" />}
                </HoverCard.Target>
                <HoverCard.Dropdown>
                  {meta.Status === "success" ? "Deployment was successful" : "Deployment is failing, check logs"}
                </HoverCard.Dropdown>
              </HoverCard>
            </Box>
          </Flex>
          
          <Collapse in={opened}>
            <Card mt={20}>
              <Anchor target="_blank" href={`https://github.com/${meta.Repo}/commit/${meta.CommitHash}`}>
                <Button leftSection={<Github size={18} />}>View on GitHub</Button>
              </Anchor>
              <Space h={15} />
              <Group wrap="nowrap" gap={4}><MessageCircle /><Text>{meta.CommitMsg}</Text></Group>
              <Group wrap="nowrap" mt={5} gap={4}><GitBranch /><Text>{meta.Branch}</Text></Group>
              <Space h={15} />
              <Timeline active={1} bulletSize={24} lineWidth={2}>
                <Timeline.Item bullet={<AlarmClockPlus size={14} />} title="Deployment Triggered">
                  <Text c="dimmed" size="xs" mt={4}>{formatDate(meta.StartAt)}</Text>
                </Timeline.Item>

                <Timeline.Item bullet={<Clock size={14} />} title="Deployment Completed">
                  <Text c="dimmed" size="xs" mt={4}>{formatDate(meta.EndAt)}</Text>
                </Timeline.Item>
              </Timeline>
            </Card>
          </Collapse>
        </Paper>
      </Center>

      <Center>
        <Paper w="90%" shadow="sm" my={20} p={20}>
          <Title order={3}>Deployment Logs</Title>
          <Space h={20} />
          <Timeline active={deploymentData?.logs.length} bulletSize={24} lineWidth={2}>
            {
              deploymentData?.logs.map(({ Title, Timestamp, ID, Message, Level }, i) => {
                let icon = <Bookmark size={16} />;
                let color = "blue";

                switch (Level) {
                  case 1:
                    color = "yellow";
                    icon = <TriangleAlert size={16} />;
                    break;
                  case 2:
                    color = "pink";
                    icon = <Bug size={16} />;
                    break;
                  case 3:
                    color = "green";
                    icon = <CheckCheck size={16} />;
                    break;
                }

                return (
                  <Timeline.Item
                    color={color}
                    bullet={icon}
                    key={(i + ID) + Timestamp}
                    title={Title}
                  >
                    <Text c="dimmed" size="sm">{Message}</Text>
                    <Text size="xs" mt={4}>{formatDate(Number(Timestamp + "000"))}</Text>
                  </Timeline.Item>
                )
              })
            }
        </Timeline>
        </Paper>
      </Center>
    </>
  )
}
