import { isFetchSuccess, showToast } from "../utils/general";
import { DeploymentDetailResp } from "../utils/types";
import { useDisclosure } from "@mantine/hooks";
import { useNavigate } from "react-router-dom";
import { formatDate } from "../utils/helpers";
import { useMeta } from "../state/use-meta";
import { useEffect, useState } from "react";
import useFetch from "../hooks/useFetch";
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

export default function Deployment() {
  const [deploymentData, setDeploymentData] = useState<DeploymentDetailResp>();
  const [opened, { toggle }] = useDisclosure(false);
  const [loading, setLoading] = useState(true);
  const { setPageTitle } = useMeta();
  const navigate = useNavigate();
  const fetchData = useFetch();

  useEffect(() => {
    const deployId = window.location.pathname.split("/").pop();

    if (!deployId || Number.isNaN(Number(deployId))) {
      showToast({
        message: "Deployment ID not found in URL",
        title: "URL Parameter Missing",
        status: "issue"
      });

      navigate("/");
      return;
    }

    setPageTitle("Deployment");

    (async () => {
      setLoading(true);
      const { signal } = new AbortController();
      const fetch = await fetchData<DeploymentDetailResp>(`/api/mg/deployment/${deployId}`, { signal });
      setLoading(false);

      if (!isFetchSuccess(fetch)) {
        navigate("/");

        return showToast({
          message: fetch.error || 'Something went wrong',
          title: 'Fetch Failed',
          status: 'issue'
        });
      }

      const { code, data } = fetch;

      if (code !== 200) {
        navigate("/");

        return showToast({
          title: 'Message From Server',
          message: data.message,
          status: 'warn'
        });
      }

      setDeploymentData(data);
    })();
  }, []);

  if (loading) return <Loading />;

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
